const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((output) => {
        Controller.updateTable(data.query, output);
      });
    });
  },

  updateTable: (query,output) => {
    var table = document.getElementById("table-body");
    var rows = `<tr><b><span style="background-color:#AED6F1">Showing total ${output["results"].length} results</span></b><tr/>`;
    tokens = output["finalQuery"]
    var rows = rows.concat(`<tr><b><i><span style="background-color:#ABEBC6">Searching for ${tokens.join(" ")}</span></i></b><tr/>`);
    for (let result of output["results"]) {
      var res = result
      for(let token of tokens) {
        var regEx = new RegExp(token, "ig");
        res = res.replaceAll(regEx, `<span style="background-color:yellow">${token}</span>`);
      }
      rows = rows.concat(`<tr><td>${res}</td><tr/>`);
    }
    table.innerHTML = rows;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
