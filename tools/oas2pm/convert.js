const Converter = require('openapi-to-postmanv2');

let input = '';
process.stdin.setEncoding('utf8');
process.stdin.on('data', (chunk) => (input += chunk));
process.stdin.on('end', () => {
    Converter.convert(
        { type: 'string', data: input },
        {
            folderStrategy: 'Tags',
            requestParametersResolution: 'Example'
        },
        (err, res) => {
            if (err) {
                console.error(err.message || String(err));
                process.exit(1);
            }
            if (!res.result) {
                console.error(res.reason || 'conversion failed');
                process.exit(2);
            }
            process.stdout.write(JSON.stringify(res.output[0].data));
        }
    );
});