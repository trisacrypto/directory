import { useToast } from '@chakra-ui/toast';
import { useState } from 'react';
import { useUpdateCertificateStep } from './useUpdateCertificateStep';
import { StepEnum } from 'types/enums';
import { validationSchema } from 'modules/dashboard/certificate/lib';
const useUploadFile = () => {
  const [isFileLoading, setIsFileLoading] = useState<boolean>(false);
  const { updateCertificateStep, wasCertificateStepUpdated, error, reset } =
    useUpdateCertificateStep();

  const toast = useToast();
  if (wasCertificateStepUpdated) {
    setIsFileLoading(false);
    reset();
    toast({
      title: 'File uploaded',
      description: 'Your file has been uploaded successfully.',
      status: 'success',
      duration: 5000,
      isClosable: true,
      position: 'top-right'
    });
  }
  if (error) {
    setIsFileLoading(false);
    reset();
    toast({
      title: 'Invalid file',
      description: error.message || 'Your json file is invalid.',
      status: 'error',
      duration: 5000,
      isClosable: true,
      position: 'top-right'
    });
  }
  const handleFileUpload = (file: any) => {
    // console.log('[handleFileUpload] file', file);
    if (file?.type !== 'application/json') {
      toast({
        title: 'Invalid file format.',
        description: `Please upload a JSON file. The maximum file size is 100KB.`,
        status: 'error',
        duration: 5000,
        isClosable: true,
        position: 'top-right'
      });
    }
    // file should'nt be up to 100kb
    if (file?.size > 100000) {
      toast({
        title: 'Invalid file size.',
        description: `JSON file size is too large. The maximum file size is 100KB. Please inspect your form or contact support for assistance.`,
        status: 'error',
        duration: 5000,
        isClosable: true,
        position: 'top-right'
      });
    }
    const reader = new FileReader();
    reader.onload = async (ev: any) => {
      const data = JSON.parse(ev.target.result);
      console.log('[handleFileUpload] data', data);
      try {
        setIsFileLoading(true);
        const basicValidationData = await validationSchema[0].validate(data, { abortEarly: true });
        const legalValidationData = await validationSchema[1].validate(data, { abortEarly: true });
        const contactValidationData = await validationSchema[2].validate(data, {
          abortEarly: true
        });
        const trisaValidationData = await validationSchema[3].validate(data, { abortEarly: true });
        const trixoValidationData = await validationSchema[4].validate(data, { abortEarly: true });

        const validationData = {
          ...basicValidationData,
          ...legalValidationData,
          ...contactValidationData,
          ...trisaValidationData,
          ...trixoValidationData
        };
        // console.log('[] validationData', validationData);
        const payload = {
          step: StepEnum.ALL,
          form: validationData?.form || validationData
        };
        // we need to create a new payload that will be used to update the certificate step

        updateCertificateStep(payload);
      } catch (e: any) {
        setIsFileLoading(false);
        console.log('[] error validationData', e.message);

        toast({
          title: 'Invalid file',
          description: e.message || 'Your json file is invalid.',
          status: 'error',
          duration: 5000,
          isClosable: true,
          position: 'top-right'
        });
      }
    };
    reader.readAsText(file);
  };

  return {
    isFileLoading,
    handleFileUpload,
    hasBeenUploaded: !!wasCertificateStepUpdated,
    hasFileUploadedFail: !!error
  };
};

export default useUploadFile;
