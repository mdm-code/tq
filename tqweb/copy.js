/**
* Copy the tq terminal command to the clipboard.
*
* @param {string} query - The query provided by the user.
* @param {string} input - The input provided by the user.
* @param {HTMLFormElement} form - The playground form with option checkboxes.
*/
const onCopyClick = async (query, input, form) => {
  const flags = [];
  if (form.querySelector('[name="tables-inline"]')?.checked) {
    flags.push('-t');
  }
  if (form.querySelector('[name="arrays-multiline"]')?.checked) {
    flags.push('-m');
  }
  if (form.querySelector('[name="indent-tables"]')?.checked) {
    flags.push('-i');
  }
  const indentSymbol = form.querySelector('[name="indent-symbol"]')?.value ?? '';
  if (indentSymbol !== '' && indentSymbol !== '  ') {
    flags.push(`-s '${indentSymbol.replace(/'/g, `'\\''`)}'`);
  }
  const flagStr = flags.length ? ` ${flags.join(' ')}` : '';
  var output = `<<EOF tq -q '${query}'${flagStr}\n${input}\nEOF`;
  try {
    navigator.clipboard.writeText(output);
  } catch (err) {
    console.log('Unable to copy the query.');
  };
};
