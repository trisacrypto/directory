import { basicDetailsValidationSchema } from 'modules/dashboard/certificate/lib/basicDetailsValidationSchema';
import { legalPersonValidationSchema } from 'modules/dashboard/certificate/lib/legalPersonValidationSchema';
import { contactsValidationSchema } from 'modules/dashboard/certificate/lib/contactsValidationSchema';
import { trisaImplementationValidationSchema } from 'modules/dashboard/certificate/lib/trisaImplementationValidationSchema';
import { trixoQuestionnaireValidationSchema } from 'modules/dashboard/certificate/lib/trixoQuestionnaireValidationSchema';

export const isBasicDetailsCompleted = async (data: any) => {
  try {
    const r = await basicDetailsValidationSchema.validate(data);
    console.log(r);
    return true;
  } catch (err) {
    return false;
  }
};

export const isLegalPersonCompleted = async (data: any) => {
  try {
    const r = await legalPersonValidationSchema.validate(data);
    console.log(r);
    return true;
  } catch (err) {
    console.log('legal person validation failed', err);
    return false;
  }
};

export const isContactsCompleted = async (data: any) => {
  try {
    const r = await contactsValidationSchema.validate(data);
    console.log(r);
    return true;
  } catch (err) {
    console.log('contacts validation failed', err);
    return false;
  }
};

export const isTrisaImplementationCompleted = async (data: any) => {
  try {
    await trisaImplementationValidationSchema.validate(data);
    return true;
  } catch (err) {
    return false;
  }
};

export const isTrixoQuestionnaireCompleted = async (data: any) => {
  try {
    await trixoQuestionnaireValidationSchema.validate(data);
    return true;
  } catch (err) {
    return false;
  }
};
