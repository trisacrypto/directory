import React from 'react';

interface StepProps {
  status: string;
  key?: number;
}
interface CertificateSteps {
  currentStep: number;
  steps: StepProps[];
  children?: React.ReactNode;
}
const initialValue: CertificateSteps = {
  currentStep: 3,
  steps: [
    {
      key: 1,
      status: 'progress'
    }
  ]
};

const CertificateStepsContext = React.createContext<
  | [Partial<CertificateSteps>, React.Dispatch<React.SetStateAction<Partial<CertificateSteps>>>]
  | typeof initialValue
>(initialValue);

const CertificateStepsProvider: React.FC<CertificateSteps> = (props) => {
  const [certificateSteps, setCertificateSteps] =
    React.useState<Partial<CertificateSteps>>(initialValue);
  return (
    <CertificateStepsContext.Provider value={[certificateSteps, setCertificateSteps]} {...props}>
      {props.children}
    </CertificateStepsContext.Provider>
  );
};

const useCertificateSteps = (): [
  Partial<CertificateSteps>,
  React.Dispatch<React.SetStateAction<Partial<CertificateSteps>>>
] => {
  const context: any = React.useContext(CertificateStepsContext);
  if (!context) {
    throw new Error('useCertificateSteps should be used within a CertificateStepsProvider');
  }

  return context;
};

export { CertificateStepsProvider, useCertificateSteps };
