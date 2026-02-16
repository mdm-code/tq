/**
* Copy the tq terminal command to the clipboard.
*
* @param {string} query - The query provided by the user.
* @param {string} input - The input provided by the user.
*/
const onCopyClick = async (query, input) => {
  var output = `<<EOF tq -q '${query}'\n${input}\nEOF`;
  try {
    navigator.clipboard.writeText(output);
  } catch (err) {
    console.log('Unable to copy the query.');
  };
};
