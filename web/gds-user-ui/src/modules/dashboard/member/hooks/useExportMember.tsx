import { useState } from "react";
import { downloadMemberToCSV } from "../utils";
const useExportMember = (member: any) => {
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const LOADING_TIMEOUT = 500;
    const exportHandler = () => {
        try{
            setIsLoading(true);
            setTimeout(() => {
                downloadMemberToCSV(member);
                setIsLoading(false);
        }, LOADING_TIMEOUT);
        } catch (error) {
            console.error('Error exporting member', error);
        }
    };
    return {
        exportHandler,
        isLoading
    };
};

export default useExportMember;
