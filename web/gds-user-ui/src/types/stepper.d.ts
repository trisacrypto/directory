// define common type
type TStep = {
  status: StepStatus;
  key?: number;
  data?: any;
};
type TPayload = {
  currentStep: number | string;
  steps: TStep[];
  lastStep: number | null;
  hasReachReviewStep?: boolean;
  hasReachSubmitStep?: boolean;
  testnetSubmitted?: boolean;
  mainnetSubmitted?: boolean;
  data?: any;
};
