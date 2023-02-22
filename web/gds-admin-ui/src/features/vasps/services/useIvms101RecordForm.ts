import { getIvms101RecordInitialValues } from '@/utils/form-references';
import { yupResolver } from '@hookform/resolvers/yup';
import React from 'react';
import { useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import * as yup from 'yup';
import { useUpdateVasp } from './update-vasp';

const validationSchema = yup.object().shape({
    national_identification: yup.object().shape({
        registration_authority: yup
            .string()
            .test('registrationAuthority', 'Registration Authority cannot be left empty', function (value) {
                if (this.parent.national_identifier_type !== 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX' && value === '') {
                    return false;
                }

                return true;
            }),
    }),
});

type FieldArrayRef = {
    addRow: () => void;
};

export default function useIvms101RecordForm({ data }: any) {
    const rhfMethods = useForm({
        defaultValues: getIvms101RecordInitialValues(data),
        resolver: yupResolver(validationSchema),
    });
    const params = useParams<{ id: string }>();

    const nameIdentifiersFieldArrayRef = React.useRef<FieldArrayRef>();
    const localNameIdentifiersFieldArrayRef = React.useRef<FieldArrayRef>();
    const phoneticNameIdentifiersFieldArrayRef = React.useRef<FieldArrayRef>();

    const nationalIdentifierType = rhfMethods.watch('national_identification.national_identifier_type');
    const isRegistrationAuthorityDisabled = React.useCallback(
        () => nationalIdentifierType === 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX',
        [nationalIdentifierType]
    );
    const { mutate: updateIvms101Record, isError, error, isLoading } = useUpdateVasp();

    const _typeahead = React.useRef<{ clear: () => void }>();

    React.useEffect(() => {
        if (nationalIdentifierType === 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX') {
            _typeahead.current?.clear();
        }
    }, [nationalIdentifierType, rhfMethods.setValue]);

    const onSubmit = async (data: any) => {
        delete data.national_identification.country_of_issue;
        delete data.national_identification.registration_authority;
        const payload = {
            entity: data,
        };

        updateIvms101Record({
            vaspId: params.id,
            data: payload,
        });
    };

    const handleAlertClose = () => {};

    const handleAddLegalNamesRow = () => {
        nameIdentifiersFieldArrayRef.current?.addRow();
    };

    const handleAddNewLocalNamesRow = () => {
        localNameIdentifiersFieldArrayRef.current?.addRow();
    };

    const handleAddNewPhoneticNamesRow = () => {
        phoneticNameIdentifiersFieldArrayRef.current?.addRow();
    };

    return {
        handleAddNewPhoneticNamesRow,
        handleAddNewLocalNamesRow,
        handleAddLegalNamesRow,
        localNameIdentifiersFieldArrayRef,
        phoneticNameIdentifiersFieldArrayRef,
        nameIdentifiersFieldArrayRef,
        isRegistrationAuthorityDisabled,
        onSubmit,
        closeAlert: handleAlertClose,
        isError,
        error,
        ...rhfMethods,
        isSubmitting: isLoading || rhfMethods.formState.isSubmitting || !rhfMethods.formState.isDirty,
        _typeahead,
    };
}
