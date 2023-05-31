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

export const getStepNumber = (step: string) => {
  switch (step) {
    case StepEnum.BASIC:
      return 1;
    case StepEnum.LEGAL:
      return 2;
    case StepEnum.CONTACTS:
      return 3;
    case StepEnum.TRISA:
      return 4;
    case StepEnum.TRIXO:
      return 5;

    default:
      return 1;
  }
};
