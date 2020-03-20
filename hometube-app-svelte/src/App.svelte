<script>
  import Files from "./Files.svelte";

  const apiBaseUrl = "http://localhost:8080/api/v1";
  const requests = {
    download: (videoUrl, videoFilename) => {
      const encodedUrl = encodeURIComponent(videoUrl);
      const encodedFilename = encodeURIComponent(videoFilename);
      const url = `${apiBaseUrl}/download?url=${encodedUrl}&filename=${encodedFilename}`;
      return fetch(url, { method: "POST" });
    }
  };

  // TODO convert to object with reset() function
  const messages = {
    error: {},
    success: {}
  };

  function download(videoUrl, videoFilename) {
    requests
      .download(videoUrl, videoFilename)
      .then(response => response.json())
      .then(json => {
        if (json.url == videoUrl) {
		  clearForm();
		  // TODO ensure videoUrl is sanitized, make it pass through URL validation
          messages.success.download = `Started downloading <a href="${videoUrl}">${videoUrl}</a> as ${videoFilename}`;
        } else {
          messages.error.download =
            "Download failed: could not start downloading";
        }
      })
      .catch(err => {
        messages.error.download = `Download failed: ${err}`;
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
    const filename = document.getElementById("videoFilename").value;
    return {
      isValid: url && filename,
      url: url,
      filename: filename
    };
  }

  function clearForm() {
    const form = document.forms[0];
    form.reset();
    form.classList.remove("was-validated");
  }

  function updateFilesList() {}

  function submit() {
    const input = validateForm();
    if (input.isValid) {
      download(input.url, input.filename);
    }

    updateFilesList();
  }
</script>

<svelte:head>
  <link
    rel="stylesheet"
    href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" />
</svelte:head>

<body class="bg-light">
  <div class="container">
    <div class="text-center">
      <h1>HomeTube</h1>
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

    <form class="needs-validation" novalidate>
      <div class="form-group">
        <label for="video-url">URL of the video to download</label>
        <input type="text" class="form-control" id="videoUrl" required />
        <div class="invalid-feedback">Please enter URL to download.</div>
      </div>
      <div class="form-group">
        <label for="video-filename">Save as... (filename for the video)</label>
        <input
          type="text"
          class="form-control"
          id="videoFilename"
          aria-describedby="videoFilenameHelp"
          required />
        <div class="invalid-feedback">
          Please enter filename to save the video.
        </div>
        <small id="videoFilenameHelp" class="form-text text-muted">
          Hopefully this will become optional in the future, taking the filename
          from the video.
        </small>
      </div>
      <button on:click={submit} type="button" class="btn btn-primary">
        Download video
      </button>
    </form>

  </div>

  <Files />

</body>
