package yaml_url_parser

import (
	"apiDocTransfer/internal/dto"
	"apiDocTransfer/pkg/config"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type YamlUrlParser struct {
	logger *zap.Logger
	cfg    config.AppConfig
}

func NewYamlUrlParser(logger *zap.Logger, cfg config.AppConfig) *YamlUrlParser {
	return &YamlUrlParser{
		logger: logger,
		cfg:    cfg,
	}
}

func (yup *YamlUrlParser) ParseYamlSources(yamlSourcesChan chan dto.YamlSourceDto) (err error) {
	f, err := os.Open(yup.cfg.SwaggerOptionsHtmlPath)
	if err != nil {
		return
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return
	}

	yup.findOptions(doc, yamlSourcesChan)

	return
}

func (yup *YamlUrlParser) findOptions(doc *goquery.Document, yamlSourcesChan chan dto.YamlSourceDto) {
	doc.Find("option").Each(func(i int, s *goquery.Selection) {
		optionValue, ok := s.Attr("value")
		if !ok {
			yup.logger.Warn("option value not found", zap.String("option", optionValue))
			return
		}

		yamlSource, processYamlDocErr := getYamlSource(yup.cfg.SwaggerDomainPrefix, optionValue)
		if processYamlDocErr != nil {
			yup.logger.Warn(
				"failed to processing yaml doc",
				zap.String("option", optionValue),
				zap.Error(processYamlDocErr),
			)
			return
		}

		yamlSourcesChan <- yamlSource
	})
}

func getYamlSource(swaggerDomainPrefix, optionValue string) (yamlSource dto.YamlSourceDto, err error) {
	yamlUrl := swaggerDomainPrefix + optionValue

	parsedUrl, err := url.Parse(yamlUrl)
	if err != nil {
		return
	}

	baseName := path.Base(parsedUrl.Path)
	ext := path.Ext(baseName)
	fileNameWithoutExt := strings.TrimSuffix(baseName, ext)

	yamlSource = dto.YamlSourceDto{
		NameWithoutExt: fileNameWithoutExt,
		Url:            yamlUrl,
	}

	return
}
