import { UseMutateFunction } from '@tanstack/react-query';

import { StepperType } from 'types/type';

export interface GetCertificateQuery {
  getCertificateStep(): void;
  certificateStep: any;
  hasCertificateStepFailed: boolean;
  wasCertificateStepFetched: boolean;
  isFetchingCertificateStep: boolean;
  error: any;
}

export interface PostCertificateMutation {
  updateCertificateStep: UseMutateFunction<PostPayloadDTO, unknown, any, unknown>;
  certificateStep: any;
  hasCertificateStepFailed: boolean;
  wasCertificateStepUpdated: boolean;
  isUpdatingCertificateStep: boolean;
  error: any;
  reset(): void;
}

export type PayloadDTO = {
  key: StepperType;
};

export type PostPayloadDTO = {
  key: StepperType;
  state: StateFormType;
  form: BasicStepType | any; // add other form types later
};
