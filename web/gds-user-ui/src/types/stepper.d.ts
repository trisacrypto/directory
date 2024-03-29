// define common type
type TStep = {
  status: StepStatus;
  key?: number;
  data?: any;
  missingFields?: any[];
  isDirty?: boolean;
};
type TDeleteStep = {
  step: TStepType;
  isDeleted: boolean;
};
type TStepType = 'basic' | 'legal' | 'contacts' | 'trisa' | 'trixo' | 'all';
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
  deletedSteps: TDeleteStep[];
};

type StateFormType = {
  current: number;
  ready_to_submit: boolean;
  started: string;
  steps: TStep[];
};

type BasicStepType = {
  organization_name: string;
  website: string;
  established_on: string;
  business_category: string;
  vasp_categories: string[];
};
