import { StepEnum } from 'types/enums';
export const getStepName = (step: number) => {
  switch (step) {
    case 1:
      return StepEnum.BASIC;
    case 2:
      return StepEnum.LEGAL;
    case 3:
      return StepEnum.CONTACTS;
    case 4:
      return StepEnum.TRISA;
    case 5:
      return StepEnum.TRIXO;

    default:
      return StepEnum.BASIC;
  }
};
