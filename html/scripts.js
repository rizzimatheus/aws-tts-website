var API_ENDPOINT = "https://ttsapi.matheusrizzi.com";

const sayButton = document.getElementById("sayButton");
sayButton.addEventListener("click", async () => {
  try {
    const inputData = {
      voice: document.querySelector("#voiceSelected option:checked").value,
      text: document.getElementById("postText").value,
    };

    const response = await fetch(API_ENDPOINT, {
      method: "POST",
      body: JSON.stringify(inputData),
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (response.ok) {
      const data = await response.json();
      document.getElementById(
        "postIDreturned"
      ).textContent = `Post ID: ${data.postId}`;
    } else {
      console.error("Error:", response.status);
    }
  } catch (error) {
    console.error("Error:", error);
  }
});

const searchButton = document.getElementById("searchButton");
searchButton.addEventListener("click", async () => {
  try {
    const postId = document.querySelector("#postId").value || "*";

    const response = await fetch(API_ENDPOINT + "?postId=" + postId);

    if (response.ok) {
      const data = await response.json();
      const postsTable = document.querySelector("#posts");

      // Remove all rows except the header
      postsTable.querySelectorAll("tr").forEach((row, index) => {
        if (index > 0) {
          row.remove();
        }
      });

      data.forEach((postData) => {
        const player = postData.url
          ? `<audio controls><source src="${postData.url}" type="audio/mpeg"></audio>`
          : "";

        const newRow = document.createElement("tr");
        newRow.innerHTML = `
                    <td>${postData.id}</td>
                    <td>${postData.voice}</td>
                    <td>${postData.text}</td>
                    <td>${postData.status}</td>
                    <td>${player}</td>
                `;

        postsTable.appendChild(newRow);
      });
    } else {
      console.error("Error:", response.status);
    }
  } catch (error) {
    console.error("Error:", error);
  }
});

const postText = document.getElementById("postText");
const charCounter = document.getElementById("charCounter");
postText.addEventListener("input", function () {
  const length = postText.value.length;
  charCounter.textContent = `Characters: ${length}`;
});
