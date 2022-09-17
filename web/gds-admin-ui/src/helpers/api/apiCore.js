import axios from 'axios';
import { defaultEndpointPrefix, getCookie } from 'utils';

const instance = axios.create({
    baseURL: defaultEndpointPrefix(),
    headers: {
        'Content-Type': 'application/json'
    }
})

// withCredentials ensures that axios fetch includes `HttpOnly` cookies in the request for CSRF protectioin
instance.defaults.withCredentials = true

const AUTH_SESSION_KEY = '__SESSION_TOKEN__';

function reauthenticate() {
    const csrfToken = getCookie('csrf_token')
    const refreshToken = getRefreshToken()
    const payload = {
        credential: refreshToken
    }

    return instance.post('/reauthenticate', payload, {
        headers: {
            'X-CSRF-TOKEN': csrfToken
        }
    })
}

let isRefreshing = false;
let subscribers = [];

function subscribeTokenRefresh(cb) {
    subscribers.push(cb);
}

function onRrefreshed(token) {
    subscribers.map((cb) => cb(token));
}

instance.interceptors.response.use(
    (response) => {
        return response;
    },
    async (error) => {
        let message;
        const originalRequest = error.config

        if (error && !error.response) {
            return Promise.reject('Network connection error')
        }

        if (error && error.response && error.response.status === 401) {
            const configData = error.response.config && error.response.config.data ? JSON.parse(error.response.config.data)?.credential : ''

            if (error.response.config.url === '/authenticate' && configData) {
                return Promise.reject(error)
            }

            if (!isRefreshing) {
                isRefreshing = true
                reauthenticate().then(response => {
                    isRefreshing = false

                    onRrefreshed(response.data)
                    subscribers = [];
                })

            }

            return new Promise((resolve) => {
                subscribeTokenRefresh(token => {
                    originalRequest.headers['Authorization'] = 'Bearer ' + token.access_token;
                    setAuthorization(token.access_token)
                    sessionStorage.setItem(AUTH_SESSION_KEY, JSON.stringify(token))

                    resolve(instance(originalRequest));
                })
            })
        }

        switch (error.response.status) {
            case 403:
                message = 'Session expired';
                sessionStorage.removeItem(AUTH_SESSION_KEY)
                setAuthorization(null)
                window.location.href = '/login'
                break;
            case 404:
                message = error || 'Sorry! the data you are looking for could not be found';
                break;
            case 500:
                message = error ?? 'Something went wrong';
                break;
            default: {
                message =
                    error.response && error.response.data ? error.response.data['error'] : error.message || error;
            }
        }
        return Promise.reject(message);
    }
);

/**
 * Sets the default authorization
 * @param {*} token
 */
const setAuthorization = (token) => {
    if (token) instance.defaults.headers.common['Authorization'] = 'Bearer ' + token;
    else delete instance.defaults.headers.common['Authorization'];
};

const setCookie = (cookie) => {
    if (cookie) {
        instance.defaults.headers.common['X-CSRF-TOKEN'] = `${cookie}`
    }
}


const getUserFromSession = () => {
    const user = sessionStorage.getItem(AUTH_SESSION_KEY);
    return user ? (typeof user == 'object' ? user : JSON.parse(user)) : null;
};

const getRefreshToken = () => {
    const user = getUserFromSession()
    return user?.refresh_token
}
class APICore {
    /**
     * Fetches data from given url
     */
    get = (url, params) => {
        let response;
        if (params) {
            var queryString = params ? params : ''
            response = instance.get(`${url}?${queryString}`);
        } else {
            response = instance.get(`${url}`, params)
        }
        return response;
    };

    getFile = (url, params) => {
        let response;
        if (params) {
            var queryString = params
                ? Object.keys(params)
                    .map((key) => key + '=' + params[key])
                    .join('&')
                : '';
            response = axios.get(`${url}?${queryString}`, { responseType: 'blob' });
        } else {
            response = axios.get(`${url}`, { responseType: 'blob' });
        }
        return response;
    };

    getMultiple = (urls, params) => {
        const reqs = [];
        let queryString = '';
        if (params) {
            queryString = params
                ? Object.keys(params)
                    .map((key) => key + '=' + params[key])
                    .join('&')
                : '';
        }

        for (const url of urls) {
            reqs.push(axios.get(`${url}?${queryString}`));
        }
        return axios.all(reqs);
    };

    /**
     * post given data to url
     */
    create = (url, data, config) => {
        return instance.post(url, data, config);
    };

    /**
     * Updates data
     */
    update = (url, data, config) => {
        return instance.put(url, data, config);
    };

    patch = (url, data, config) => {
        return instance.patch(url, data, config)
    }

    /**
     * Deletes data
     */
    delete = (url, params) => {
        return instance.delete(url, params);
    };

    /**
     * post given data to url with file
     */
    createWithFile = (url, data) => {
        const formData = new FormData();
        for (const k in data) {
            formData.append(k, data[k]);
        }

        const config = {
            headers: {
                ...axios.defaults.headers,
                'content-type': 'multipart/form-data',
            },
        };
        return axios.post(url, formData, config);
    };

    /**
     * post given data to url with file
     */
    updateWithFile = (url, data) => {
        const formData = new FormData();
        for (const k in data) {
            formData.append(k, data[k]);
        }

        const config = {
            headers: {
                ...axios.defaults.headers,
                'content-type': 'multipart/form-data',
            },
        };
        return axios.patch(url, formData, config);
    };

    setLoggedInUser = (session) => {
        if (session) sessionStorage.setItem(AUTH_SESSION_KEY, JSON.stringify(session));
        else {
            sessionStorage.removeItem(AUTH_SESSION_KEY);
        }
    };

    /**
     * Returns the logged in user
     */
    getLoggedInUser = () => {
        return getUserFromSession();
    };

    deleteUserSession = () => {
        this.setLoggedInUser(null)
        setAuthorization(null)
        window.location.href = '/login'
    }

    getRefreshSessionToken = () => {
        return getRefreshToken()
    }
}

/*
Check if token available in session
*/
const user = getUserFromSession();
if (user && user.access_token) {
    setAuthorization(user.access_token);
}

export { APICore, setAuthorization, setCookie };
