import { useEffect, useState } from "react";
import { getSubmissionStatus } from "../service";
import { upperCaseFirstLetter } from "utils/utils";

const useSubmissionStatus = () => {
  const [error, setError] = useState<any>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [status, setStatus] = useState<any>(null);

    const fetchSubmissionStatus = async () => {
      setIsLoading(true);
      try {
        const response = await getSubmissionStatus();
        if (!response) setError('No data found');
        setStatus(response);
      } catch (e: any) {
        if (!e?.data?.success) {
          setError(upperCaseFirstLetter(e?.data?.error));
        } else {
          setError('Something went wrong.');
        }
      } finally {
        setIsLoading(false);
      }
    };

    useEffect(() => {
      fetchSubmissionStatus();
    }, []);

    return { status, isLoading, error };
  };

export default useSubmissionStatus;
