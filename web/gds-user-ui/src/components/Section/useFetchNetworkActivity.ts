import { useState } from "react";
import { networkActivity } from "./service";

const useFetchNetworkActivity = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [data, setData] = useState<any>(null);
    const [error, setError] = useState<any>(null);

    const fetchNetworkActivity = async () => {
        setIsLoading(true);
        try {
            const response = await networkActivity();
            setData(response);
        } catch (e: any) {
            if (!e?.response?.data?.success) {
                setError(e?.response?.data?.error);
            } else {
                setError("Something went wrong.");
            }
        }
        setIsLoading(false);
    };
    fetchNetworkActivity();
    return { data, isLoading, error };
};

export default useFetchNetworkActivity;
