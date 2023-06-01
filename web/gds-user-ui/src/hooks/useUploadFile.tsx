import { useToast } from '@chakra-ui/toast';
import { useState } from 'react';
import { useUpdateCertificateStep } from './useUpdateCertificateStep';
import { StepEnum } from 'types/enums';
import { validationSchema } from 'modules/dashboard/certificate/lib';
import { handleError } from 'utils/utils';
const useUploadFile = () => {
  console.log('[useUploadFile] init');
  const { updateCertificateStep, wasCertificateStepUpdated, error, reset } =
    useUpdateCertificateStep();

  const toast = useToast();
  if (wasCertificateStepUpdated) {
    reset();
    toast({
      title: 'File uploaded',
      description: 'Your file has been uploaded successfully',
      status: 'success',
      duration: 5000,
      isClosable: true,
      position: 'top-right'
    });
  }
  if (error) {
    reset();
    toast({
      title: 'Invalid file',
      description: error.message || 'Your json file is invalid',
      status: 'error',
      duration: 5000,
      isClosable: true,
      position: 'top-right'
    });
  }
  const [isFileLoading, setIsFileLoading] = useState<boolean>(false);
  const handleFileUpload = (file: any) => {
    console.log('[handleFileUpload] file', file);
    setIsFileLoading(true);
    const reader = new FileReader();
    reader.onload = async (ev: any) => {
      const data = JSON.parse(ev.target.result);
      try {
        const validationData = await validationSchema[0].validate(data, { abortEarly: false });
        console.log('[] validationData', validationData);
        const payload = {
          step: StepEnum.ALL,
          form: validationData?.form || validationData
        };
        updateCertificateStep(payload);
      } catch (e: any) {
        console.log('[] error validationData', e);

        handleError(e, `[Invalid file], it's missing some required fields : ${e.message}`);
      } finally {
        setIsFileLoading(false);
      }
    };
    reader.readAsText(file);
  };

  return { isFileLoading, handleFileUpload, hasBeenUploaded: wasCertificateStepUpdated };
};

export default useUploadFile;