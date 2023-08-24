import { useState } from "react";
import { networkActivity } from "./service";

const useFetchNetworkActivity = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [data, setData] = useState<any>(null);
    const [error, setError] = useState<any>(null);
    
    const fetchNetworkActivity = async () => {
        setIsLoading(true);
        try {
            const res = await networkActivity();
            if (!res.mainnet && !res.testnet) setError("No network activity found.");
            setData(res);
        } catch (e: any) {
            if (!e?.res?.data?.success) {
                setError(e?.res?.data?.error);
            } else {
                setError("Something went wrong.");
            }
        } finally{
            setIsLoading(false);
        }
    };
    fetchNetworkActivity();
    return { data, isLoading, error };
};

export default useFetchNetworkActivity;
