import axiosInstance from "utils/axios";

export const networkActivity = async () => {
    const response = await axiosInstance.get(`/network/activity`);
    return response.data;
};
