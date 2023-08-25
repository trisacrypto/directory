import axiosInstance from "utils/axios";

export const networkActivity = async () => {
    const res = await axiosInstance.get(`/network/activity`);
    return res.data;
};
