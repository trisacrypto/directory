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
  updatedCertificateStep: any;
  hasCertificateStepFailed: boolean;
  wasCertificateStepUpdated: boolean;
  isUpdatingCertificateStep: boolean;
  error: any;
  reset(): void;
}

export interface DeleteCertificateMutation {
  deleteCertificateStep: UseMutateFunction<Pick<PostPayloadDTO, 'step'>, unknown, any, unknown>;
  deletedCertificateStep: any;
  hasCertificateStepFailed: boolean;
  wasCertificateStepDeleted: boolean;
  isDeletingCertificateStep: boolean;
  error: any;
  reset(): void;
}

export type PayloadDTO = {
  key: StepperType;
};

export type PostPayloadDTO = {
  step: StepperType;
  form: any; // add other form types later
};
