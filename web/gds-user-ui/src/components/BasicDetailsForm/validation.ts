import * as yup from 'yup';

export const ValidationSchema = yup.object().shape({
  website: yup.string().url().required(),
  business_category: yup.string(),
  vasp_categories: yup.array().of(yup.string()),
  established_on: yup.date()
});

export const getDefaultValue = () => {
  return {
    website: '',
    business_category: '',
    vasp_categories: [],
    established_on: ''
  };
};
