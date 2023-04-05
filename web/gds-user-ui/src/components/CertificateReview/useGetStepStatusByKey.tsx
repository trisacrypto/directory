import { useSelector } from 'react-redux';
import { getSteps } from 'application/store/selectors/stepper';

const useGetStepStatusByKey = (key?: number) => {
  const k = key || 1;
  const steps = useSelector(getSteps);
  const hasErrorField = steps[k]?.status === 'error';
  const isCompleted = steps[k]?.status === 'complete';
  const isPending = steps[k]?.status === 'progress';
  const requiredMissingFields = key && steps?.map((step: any) => step?.status === 'error');

  return {
    hasErrorField,
    isCompleted,
    isPending,
    requiredMissingFields
  };
};

export default useGetStepStatusByKey;
