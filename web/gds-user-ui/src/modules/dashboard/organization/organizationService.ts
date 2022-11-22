
import axiosInstance from 'utils/axios';


export const GetAllOrganisations = async () => {
    const response = await axiosInstance.get(`/organisations`);
    return response;
};

export const GetOrganisation = async (id: string) => {
    const response = await axiosInstance.get(`/organisations/${id}`);
    return response;
};

export const CreateOrganisation = async (data: any) => {
    const response = await axiosInstance.post(`/organisations`, data);
    return response;
};

export const UpdateOrganisation = async (id: string, data: any) => {
    const response = await axiosInstance.put(`/organisations/${id}`, data);
    return response;
};
