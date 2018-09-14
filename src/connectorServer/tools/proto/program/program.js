/**
 * Created by yangsong on 16/1/23.
 */
var Program = require('commander');

Program
    .option('-p, --protoFile <n>', 'proto file')
    .option('-n, --projectName <n>', 'project name')
    .parse(process.argv);

module.exports = Program;
