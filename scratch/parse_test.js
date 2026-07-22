
const parser = require('@babel/parser');
const fs = require('fs');
const code = fs.readFileSync('scratch/test_cw.js', 'utf8');
try {
  parser.parse(code, { sourceType: 'module', plugins: ['jsx'] });
  console.log('Babel parsed successfully!');
} catch (e) {
  console.log('Babel error:', e.message, 'at line', e.loc.line, 'col', e.loc.column);
}
