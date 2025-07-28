let i = 2;
document.getElementById('add-new-alert').onclick = addNewAlertHandler;

document.getElementById('remove-alert').onclick = function () {
    let container = document.getElementById('alert-fields');
    const removeBtn = document.getElementById('remove-alert');
    if (container.children.length > 1) {
      container.removeChild(container.lastElementChild);
      i--;
      if (i < 3) {
        const addBtn = document.getElementById('add-new-alert');
        addBtn.classList.remove('disabled');
        addBtn.innerText = "Add New Alert";
        addBtn.onclick = addNewAlertHandler;
      }
      if (container.children.length === 1) {
        removeBtn.classList.add('disabled');
        removeBtn.onclick = null;
      }
    }
};

// Enable remove button if more than one alert exists on page load
window.onload = function() {
  const container = document.getElementById('alert-fields');
  const removeBtn = document.getElementById('remove-alert');
  if (container.children.length > 1) {
    removeBtn.classList.remove('disabled');
    removeBtn.innerText = "Remove Alert";
    removeBtn.onclick = document.getElementById('remove-alert').onclick;
  } else if (container.children.length < 2) {
    removeBtn.classList.add('disabled');
    removeBtn.innerText = "Minimum of 1 alert required";
    removeBtn.onclick = null;
  }
};

function addNewAlertHandler() {
  let template = `
      <legend >Alert ${i}</legend>
      <div class="form-group">
        <label for="ChannelName">Channel Name</label>
        <input type="text" class="form-control" id="ChannelName" name="ChannelName" 
          placeholder="Enter Slack channel name"
          maxlength="20" required>
        <small class="form-text text-muted"><strong>Please ensure you don't include the # at the beginning e.g example-alerts</strong></small>
      </div>
      <div class="form-group">
        <label for="SlackWebhookURL">Slack Webhook URL</label>
        <input type="password" class="form-control" id="SlackWebhookURL" name="SlackWebhookURL"
          placeholder="Enter Slack webhook URL" 
          pattern="https://hooks.slack.com/services/[A-Za-z0-9]+/[A-Za-z0-9]+/[A-Za-z0-9]+"
          required>
        <input type="checkbox" id="showWebhookURL" name="showWebhookURL" onclick="ShowWebhookURL(this)">Show URL
      </div>
      <div class="form-group">
        <label for="Severity">Severity</label>
        <input type="text" class="form-control" id="Severity" name="Severity"
          placeholder="Enter severity"
          maxlength="10" required>
        <small class="form-text text-muted"><strong>e.g critical, warning, info </strong></small>
      </div>`;

    let container = document.getElementById('alert-fields');
    let div = document.createElement('fieldset');
    div.innerHTML = template;
    container.appendChild(div);

    i++;
    if (i > 3) {
      const addBtn = document.getElementById('add-new-alert');
      addBtn.classList.add('disabled');
      addBtn.innerText = "Maximum of 3 alerts reached";
      addBtn.onclick = null;
    }
    // Enable remove button when more than one alert exists
    const removeBtn = document.getElementById('remove-alert');
    if (container.children.length > 1) {
        removeBtn.classList.remove('disabled');
        removeBtn.innerText = "Remove Alert";
        // DON'T redefine onclick here - it's already defined at the top level
    }

    document.getElementById('remove-alert').onclick = function () {
    let container = document.getElementById('alert-fields');
    if (container.children.length > 1) {
      container.removeChild(container.lastElementChild);
      i--;
      if (i < 3) {
        const addBtn = document.getElementById('add-new-alert');
        addBtn.classList.remove('disabled');
        addBtn.innerText = "Add New Alert";
        addBtn.onclick = addNewAlertHandler;
      }
    }
  }
}

