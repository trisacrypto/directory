import { reviewValidationSchema } from './reviewValidationSchema';
import { legalPersonValidationSchema } from './legalPersonValidationSchema';
import { basicDetailsValidationSchema } from './basicDetailsValidationSchema';
import { contactsValidationSchema } from './contactsValidationSchema';
import { trisaImplementationValidationSchema } from './trisaImplementationValidationSchema';
import { trixoQuestionnaireValidationSchema } from './trixoQuestionnaireValidationSchema';

export const validationSchema = [
  basicDetailsValidationSchema,
  legalPersonValidationSchema,
  contactsValidationSchema,
  trisaImplementationValidationSchema,
  trixoQuestionnaireValidationSchema,
  reviewValidationSchema
];

export const allValidationSchema = validationSchema.reduce((acc, schema) => {
  return acc.concat(schema.fields);
}, []);
