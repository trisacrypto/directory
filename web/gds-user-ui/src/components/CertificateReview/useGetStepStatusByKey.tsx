import { useSelector } from 'react-redux';
import { getSteps } from 'application/store/selectors/stepper';

const useGetStepStatusByKey = (key = 1) => {
  const k = key - 1;
  const steps = useSelector(getSteps);
  const hasErrorField = steps[k - 1]?.status === 'error';
  const isCompleted = steps[k]?.status === 'complete';
  const isPending = steps[k]?.status === 'progress';
  const requiredMissingFields = steps?.map((step: any) => step?.status === 'progress');

  return {
    hasErrorField,
    isCompleted,
    isPending,
    requiredMissingFields
  };
};

export default useGetStepStatusByKey;
