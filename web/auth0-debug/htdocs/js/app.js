let auth0 = null;

const fetchAuthConfig = () => fetch("/auth_config.json");

const configureClient = async () => {
  const response = await fetchAuthConfig();
  const config = await response.json();

  auth0 = await createAuth0Client({
    domain: config.domain,
    client_id: config.clientId,
    audience: config.audience
  });
};

const login = async () => {
  await auth0.loginWithRedirect({
    redirect_uri: window.location.origin
  });
};

const logout = async() => {
  auth0.logout({
    returnTo: window.location.origin
  });
};

const updateUI = async () => {
  const isAuthenticated = await auth0.isAuthenticated();

  document.getElementById("btn-logout").disabled = !isAuthenticated;
  document.getElementById("btn-login").disabled = isAuthenticated;

  if (isAuthenticated) {
    let user = await auth0.getUser();

    document.getElementById("gated-content").classList.remove("invisible");
    document.getElementById("img-whoami").src = user.picture;
    document.getElementById("app-banner").textContent = "Welcome " + user.name;
    document.getElementById("ipt-access-token").innerHTML = await auth0.getTokenSilently();

    console.log(await auth0.getTokenSilently());
    document.getElementById("ipt-user-profile").textContent = JSON.stringify(user);
  } else {
    document.getElementById("gated-content").classList.add("invisible");
    document.getElementById("img-whoami").src = "/img/logo.svg";
    document.getElementById("app-banner").textContent = "Please Sign In";
  }

};

window.onload = async() => {
  await configureClient();
  updateUI();

  const isAuthenticated = await auth0.isAuthenticated();
  if (isAuthenticated) {
    // show the gated content without parsing the query params
    return;
  }

  // Check the code and state parameters
  const query = window.location.search;
  if (query.includes("code=") && query.includes("state=")) {
    // Process the login state
    await auth0.handleRedirectCallback();

    updateUI();

    // Use replaceState to redirect the user away and remove querystring parameters
    window.history.replaceState({}, document.title, "/");
  }

};
