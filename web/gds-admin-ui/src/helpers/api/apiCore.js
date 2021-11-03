import axios from 'axios';
import { defaultEndpointPrefix, getCookie } from '../../utils';
import jwtDecode from 'jwt-decode'
import toast from 'react-hot-toast';

axios.defaults.headers.post['Content-Type'] = 'application/json';
axios.defaults.baseURL = defaultEndpointPrefix();

// withCredentials ensures that axios fetch includes `HttpOnly` cookies in the request for CSRF protectioin
axios.defaults.withCredentials = true

const AUTH_SESSION_KEY = '__SESSION_TOKEN__';

axios.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        let message;

        if (error && !error.response) {
            return Promise.reject('Network connection error')
        }

        if (error && error.response && error.response.status === 404) {
            // window.location.href = '/not-found';
        } else if (error && error.response && error.response.status === 403) {
            sessionStorage.removeItem(AUTH_SESSION_KEY)
            window.location.href = '/login';
        } else {
            switch (error.response.status) {
                case 401:
                    message = 'Invalid credentials';
                    window.location.href = '/login'
                    break;
                case 403:
                    message = 'Access Forbidden';
                    break;
                case 404:
                    message = 'Sorry! the data you are looking for could not be found';
                    break;
                case 500:
                    message = 'Something went wrong';
                    break;
                default: {
                    message =
                        error.response && error.response.data ? error.response.data['message'] : error.message || error;
                }
            }
            return Promise.reject(message);
        }
    }
);

/**
 * Sets the default authorization
 * @param {*} token
 */
const setAuthorization = (token) => {
    if (token) axios.defaults.headers.common['Authorization'] = 'Bearer ' + token;
    else delete axios.defaults.headers.common['Authorization'];
};

const setCookie = (cookie) => {
    if (cookie) {
        axios.defaults.headers.common['X-CSRF-TOKEN'] = `${cookie}`
    }
}

const isValidRefreshToken = (token) => {
    if (token) {
        const decoded = jwtDecode(token)
        const currentTime = Date.now() / 1000
        return currentTime > decoded.nbf && currentTime < decoded.exp

    }

    return false;
}


const getUserFromSession = () => {
    const user = sessionStorage.getItem(AUTH_SESSION_KEY);
    return user ? (typeof user == 'object' ? user : JSON.parse(user)) : null;
};
class APICore {
    /**
     * Fetches data from given url
     */
    get = (url, params) => {
        let response;
        if (params) {
            var queryString = params
                ? Object.keys(params)
                    .map((key) => key + '=' + params[key])
                    .join('&')
                : '';
            response = axios.get(`${url}?${queryString}`, params);
        } else {
            response = axios.get(`${url}`, params)
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
        return axios.post(url, data, config);
    };

    /**
     * Updates patch data
     */
    updatePatch = (url, data) => {
        return axios.patch(url, data);
    };

    /**
     * Updates data
     */
    update = (url, data) => {
        return axios.put(url, data);
    };

    /**
     * Deletes data
     */
    delete = (url) => {
        return axios.delete(url);
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

    reauthenticate = (payload) => {
        const cookie = getCookie('csrf_token')

        axios.post('/reauthenticate', payload, {
            headers: {
                'X-CSRF-TOKEN': cookie
            }
        }).then(res => {
            this.setLoggedInUser(res.data)
            setAuthorization(res.data.access_token)
            return true
        }).catch((err) => {
            console.log('[Error]', err)
            toast.error("Something wrong happen while refreshing the token")
            this.setLoggedInUser(null)
            setAuthorization(null)
            return false
        })

        return null
    }

    isUserAuthenticated = () => {
        const user = this.getLoggedInUser();
        if (!user) {
            return false;
        }

        const decodedAccessToken = jwtDecode(user.access_token);
        const currentTime = Date.now() / 1000;

        const payload = {
            credential: user.refresh_token
        }

        if (currentTime < decodedAccessToken.exp && currentTime > decodedAccessToken.nbf) {
            // The access token is valid -- we could just return true here
            // Alternatively, we could check if we're in that small window of time where we can reauthenticate when the access token is valid:
            if (isValidRefreshToken(user.refresh_token)) {
                this.reauthenticate(payload)
            }
            return true;
        } else {
            // access token is invalid, check if we can reauthenticate
            if (isValidRefreshToken(user.refresh_token)) {
                this.reauthenticate(payload)
            }
            // neither the access nor the refresh token is valid any longer
            this.setLoggedInUser(null)
            setAuthorization(null)
            return false;
        }
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

    setUserInSession = (modifiedUser) => {
        let userInfo = sessionStorage.getItem(AUTH_SESSION_KEY);
        if (userInfo) {
            const { token, user } = JSON.parse(userInfo);
            this.setLoggedInUser({ token, ...user, ...modifiedUser });
        }
    };

    deleteUserSession = () => {
        this.setLoggedInUser(null)
        setAuthorization(null)
        window.location.href = '/login'
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
