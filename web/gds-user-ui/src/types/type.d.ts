// define common type
type TStep = {
  status: StepStatus;
  key?: number;
  data?: any;
};

type StepStatus = 'complete' | 'progress' | 'incomplete';
