package parser

import (
	"apiDocTransfer/internal/dto"
	"apiDocTransfer/internal/service/open_api_converter"
	"apiDocTransfer/internal/service/postman_client"
	"apiDocTransfer/internal/service/yaml_doc_processor"
	"apiDocTransfer/internal/service/yaml_url_parser"
	"apiDocTransfer/pkg/config"
	"context"
	"fmt"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

const numWorkers = 5

type Parser struct {
	cfg           config.AppConfig
	logger        *zap.Logger
	client        *http.Client
	postmanClient *postman_client.PostmanClient
}

type folderResult struct {
	name       string
	collection map[string]any
}

func NewParser(cfg config.AppConfig, logger *zap.Logger, client *http.Client, postmanClient *postman_client.PostmanClient) *Parser {
	return &Parser{
		cfg:           cfg,
		logger:        logger,
		client:        client,
		postmanClient: postmanClient,
	}
}

func (p *Parser) Parse(ctx context.Context) error {
	yamlUrlParser := yaml_url_parser.NewYamlUrlParser(p.logger, p.cfg)

	sources := make(chan dto.YamlSourceDto)
	results := make(chan folderResult)

	var produceErr error
	go func() {
		defer close(sources)
		parseSourceErr := yamlUrlParser.ParseYamlSources(sources)
		if parseSourceErr != nil {
			produceErr = fmt.Errorf("failed to parse yaml sources: %w", parseSourceErr)
		}
	}()

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Go(func() {
			p.processSources(ctx, sources, results)
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if err := p.postmanClient.UpsertFolder(ctx, p.cfg.PostmanCollectionID, res.name, res.collection); err != nil {
			p.logger.Error("postmanClient.UpsertFolder error", zap.Error(err), zap.String("folder", res.name))
		}

		p.logger.Info("postmanClient.UpsertFolder success", zap.String("folder", res.name))
	}

	return produceErr
}

func (p *Parser) processSources(ctx context.Context, sources chan dto.YamlSourceDto, results chan folderResult) {
	for source := range sources {
		if res, ok := p.processYamlSource(ctx, source); ok {
			results <- res
		}
	}
}

func (p *Parser) processYamlSource(ctx context.Context, yamlSource dto.YamlSourceDto) (folderResult, bool) {
	yamlDoc, err := yaml_doc_processor.GetYamlDoc(p.client, yamlSource.Url, yamlSource.NameWithoutExt, p.cfg.SwaggerAuthToken)
	if err != nil {
		p.logger.Error("get yaml doc error", zap.Error(err), zap.String("url", yamlSource.Url))
		return folderResult{}, false
	}

	collection, err := open_api_converter.ConvertOpenAPI(ctx, p.cfg.Oas2pmPath, yamlDoc)
	if err != nil {
		p.logger.Error("convert openapi error", zap.Error(err), zap.Any("source", yamlSource))
		return folderResult{}, false
	}

	result := folderResult{
		name:       yamlSource.NameWithoutExt,
		collection: collection,
	}

	return result, true
}
