<script>
  import Items from "./Items.svelte";
  import { items } from "./stores.js";
  import { onMount } from "svelte";

  export let apiBaseUrl;

  const requests = {
    download: videoUrl => {
      const encodedUrl = encodeURIComponent(videoUrl);
      const url = `${apiBaseUrl}/download?url=${encodedUrl}`;
      return fetch(url, { method: "POST" });
    },
    listDownloaded: () => {
      const url = `${apiBaseUrl}/list/downloaded`;
      return fetch(url, { method: "GET" });
    }
  };

  // TODO convert to object with reset() function
  const messages = {
    error: {},
    success: {}
  };

  function download(videoUrl) {
    requests
      .download(videoUrl)
      .then(response => response.json())
      .then(json => {
        if (json.url == videoUrl) {
          clearForm();
          // TODO ensure videoUrl is sanitized, make it pass through URL validation
          messages.success.download = `Started downloading <a href="${videoUrl}">${videoUrl}</a>`;
        } else {
          messages.error.download =
            "Download failed: could not start downloading";
        }
      })
      .catch(err => {
        messages.error.download = `Download failed: ${err}`;
      });
  }

  function updateFilesList() {
    console.log("Updating list of downloaded files...");
    requests
      .listDownloaded()
      .then(response => response.json())
      .then(json => {
        if (json.items) {
          items.update(prev => json.items);
        } else {
          messages.error.listDownloaded =
            "Could not get list of downloaded files";
        }
      })
      .catch(err => {
        messages.error.listDownloaded = `Could not get list of downloaded files: ${err}`;
      });
  }

  function validateForm() {
    messages.error = {};
    messages.success = {};

    const form = document.forms[0];
    if (form.checkValidity() === false) {
      event.preventDefault();
      event.stopPropagation();
    }
    form.classList.add("was-validated");

    const url = document.getElementById("videoUrl").value;
    return {
      isValid: !!url,
      url: url
    };
  }

  function clearForm() {
    const form = document.forms[0];
    form.reset();
    form.classList.remove("was-validated");
  }

  function submit() {
    const input = validateForm();
    if (input.isValid) {
      download(input.url);
    }

    updateFilesList();
  }

  onMount(updateFilesList);
</script>

<svelte:head>
  <link
    rel="stylesheet"
    href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" />
</svelte:head>

<body class="bg-light">
  <div class="container">
    <div class="text-center">
      <h1>Hometube</h1>
      <p class="lead">
        Download videos from YouTube or other streaming sites using the
        <code>youtube-dl</code>
        tool.
      </p>
    </div>

    {#if messages.error.download}
      <div class="alert alert-danger" role="alert">
        {messages.error.download}
      </div>
    {/if}

    {#if messages.success.download}
      <div class="alert alert-success" role="alert">
        {@html messages.success.download}
      </div>
    {/if}

    <form on:submit|preventDefault={submit} class="needs-validation" novalidate>
      <div class="form-group">
        <label for="video-url">URL of the video or playlist to download</label>
        <input type="text" class="form-control" id="videoUrl" required />
        <div class="invalid-feedback">Please enter URL to download.</div>
      </div>
      <button type="submit" class="btn btn-primary">Download</button>
    </form>

  </div>

  <Items />

</body>
