// define common type
type TStep = {
  status: StepStatus;
  key?: number;
  data?: any;
};
type TStepStatus = 'progress' | 'success' | 'error';
type TStatusKey = 'testnet' | 'mainnet';
type TPayload = {
  currentStep: number | string;
  steps: TStep[];
  lastStep: number | null;
  hasReachReviewStep?: boolean;
  hasReachSubmitStep?: boolean;
  testnetSubmitted?: boolean;
  mainnetSubmitted?: boolean;
  status?: Record<TStatusKey, StepStatus>;
  data?: any;
};
