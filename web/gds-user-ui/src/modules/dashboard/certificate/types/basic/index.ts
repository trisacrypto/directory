export type PostBasicPayloadDTO = {
  state: StateFormType;
  form: BasicStepType;
};

export type BasicStepType = {
  organization_name: string;
  website: string;
  established_on: string;
  business_category: string;
  vasp_categories: string[];
};
