import { StepStatus } from 'types/type';
import { NotSavedSectionStatus, SavedSectionStatus } from '.';

type SectionStatusProps = {
  status: StepStatus;
};

const SectionStatus: React.FC<SectionStatusProps> = ({ status }) => {
  return (
    <>
      {status === 'complete' || status === 'incomplete' ? (
        <SavedSectionStatus />
      ) : (
        <NotSavedSectionStatus />
      )}
    </>
  );
};

export { SectionStatus };
