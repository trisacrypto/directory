import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

const _i18n = setupI18n();

export const addCollaboratorFormValidationSchema = yup
    .object()
    .shape({
        email: yup
            .string()
            .email(_i18n._(t`Email is not valid.`))
            .required(_i18n._(t`Email is required.`)),
        name: yup.string().required(),
        agreed: yup
            .boolean()
            .oneOf([true], _i18n._(t`You must agree to the terms and conditions`))
            .default(false)
    }).required();


export const ADD_COLLABORATOR_FORM_METHOD = {
    resolver: yupResolver(addCollaboratorFormValidationSchema),
    defaultValues: {
        email: '',
        name: '',
        agreed: false
    }
};

