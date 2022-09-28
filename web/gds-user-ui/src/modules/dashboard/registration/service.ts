import axiosInstance, { setAuthorization } from 'utils/axios';
export const getRegistrationData = async () => {
  setAuthorization();
  const response = await axiosInstance.get(`/register`);
  return response;
};
export const postRegistrationData = async (data: any) => {
  setAuthorization();
  const response = await axiosInstance.put(`/register`, { ...data });
  return response;
};

// submit testnet registration endpoint request

export const submitTestnetRegistration = async () => {
  setAuthorization();
  const response = await axiosInstance.post(`/register/testnet`);
  return response;
};

// submit mainnet registration endpoint request

export const submitMainnetRegistration = async () => {
  setAuthorization();
  const response = await axiosInstance.post(`/register/mainnet`);
  return response;
};

// set default state for registration

export const setRegistrationDefaultState = async () => {
  setAuthorization();
  const response = await axiosInstance.put(`/register`, {
    state: {
      current: 1,
      ready_to_submit: false,
      steps: [
        {
          key: 1,
          status: 'progress'
        }
      ]
    }
  });
  return response;
};

// get submission status

export const getSubmissionStatus = async () => {
  setAuthorization();
  const response = await axiosInstance.get(`/registration`);
  return response;
};
