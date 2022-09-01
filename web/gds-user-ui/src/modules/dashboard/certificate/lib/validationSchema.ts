import { reviewValidationSchema } from './reviewValidationSchema';
import { legalPersonValidationSchemam } from './legalPersonValidationSchema';
import { basicDetailsValidationSchema } from './basicDetailsValidationSchema';
import { contactsValidationSchema } from './contactsValidationSchema';
import { trisaImplementationValidationSchema } from './trisaImplementationValidationSchema';
import { trixoQuestionnaireValidationSchema } from './trixoQuestionnaireValidationSchema';

export const validationSchema = [
  basicDetailsValidationSchema,
  legalPersonValidationSchemam,
  contactsValidationSchema,
  trisaImplementationValidationSchema,
  trixoQuestionnaireValidationSchema,
  reviewValidationSchema
];
