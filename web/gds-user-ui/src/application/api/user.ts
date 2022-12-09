import axiosInstance from "utils/axios";

export const updateUserFullName = async (fullName: string) => {
    return await axiosInstance.patch("/users", {
        name: fullName
    });
};
