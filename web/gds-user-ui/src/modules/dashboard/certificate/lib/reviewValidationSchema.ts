import * as yup from 'yup';

export const reviewValidationSchema = yup
  .object()
  .shape({
    state: yup.object().shape({
      current: yup.number(),
      steps: yup.array().of(
        yup.object().shape({
          status: yup.string().oneOf(['complete', 'incomplete', 'pending']),
          key: yup.number().required()
        })
      ),
      reach_submit_step: yup.boolean().default(false)
    })
  })
  .notRequired();
