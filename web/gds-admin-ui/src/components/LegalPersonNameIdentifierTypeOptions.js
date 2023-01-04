import { NAME_IDENTIFIER_TYPE } from '@/constants';

function LegalPersonNameIdentifierTypeOptions() {
  return (
    <>
      {Object.entries(NAME_IDENTIFIER_TYPE).map(([k, v]) => (
        <option key={k} value={k}>{`${v} Name`}</option>
      ))}
    </>
  );
}

export default LegalPersonNameIdentifierTypeOptions;
