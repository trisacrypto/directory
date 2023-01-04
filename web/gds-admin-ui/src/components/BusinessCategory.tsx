import { BUSINESS_CATEGORY } from '@/constants';

function BusinessCategory() {
  return (
    <>
      <option value="" />
      {Object.entries(BUSINESS_CATEGORY).map(([k, v]) => (
        <option value={k} key={k}>
          {v}
        </option>
      ))}
    </>
  );
}

export default BusinessCategory;
