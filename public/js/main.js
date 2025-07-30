const baseUrl = `${window.location.protocol}//${window.location.host}`;
// Search form handler
document.getElementById("searchForm").addEventListener("submit", function (e) {
  e.preventDefault();
  const query = document.getElementById("query").value;
  const results = document.getElementById("searchResults");

  results.innerHTML = "<p>Sending request...</p>";
  fetch(`${baseUrl}/api/v1/search`, {
    method: "POST",
    body: JSON.stringify({
      query,
    }),
    headers: {
      "Content-type": "application/json",
    },
  })
    .then((res) => {
      if (!res.ok) return (results.innerHTML = "Some Error");
      return res.json();
    })
    .then((data) => {
      results.innerHTML = `
              <h4>Response:</h4>
               <pre style="background-color: #f6f8fa; padding: 10px; border-radius: 5px;">
              ${JSON.stringify(data)}
              </pre>
              `;
    })
    .catch((err) => {
      console.log(err);
    });
});

// Test API handler
document.getElementById("testApi").addEventListener("click", function () {
  const results = document.getElementById("testResults");

  results.innerHTML = "<p>Testing API...</p>";
  fetch(`${baseUrl}/api/v1/test`)
    .then((res) => {
      if (!res.ok) return (results.innerHTML = "Some Error");
      return res.text();
    })
    .then((text) => {
      results.innerHTML = `
              <h4>Response:</h4>
               <pre style="background-color: #f6f8fa; padding: 10px; border-radius: 5px;">${text}</pre>
              `;
    })
    .catch((err) => {
      console.log(err);
    });
});
