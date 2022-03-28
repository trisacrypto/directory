import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_TRISA_BASE_URL,
  headers: {
    accept: 'application/json'
  }
});

export default axiosInstance;
