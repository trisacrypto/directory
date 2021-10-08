import React from 'react';
import { Trans } from "@lingui/macro"


const TransliterationMethodCode = () => {
  return (
    <>
    <option value={1}><Trans>Arabic (Arabic language)</Trans></option>
    <option value={2}><Trans>Arabic (Persian language)</Trans></option>
    <option value={3}><Trans>Armenian</Trans></option>
    <option value={4}><Trans>Cyrillic</Trans></option>
    <option value={5}><Trans>Devanagari & related Indic</Trans></option>
    <option value={6}><Trans>Georgian</Trans></option>
    <option value={7}><Trans>Greek</Trans></option>
    <option value={8}><Trans>Han (Hanzi, Kanji, Hanja)</Trans></option>
    <option value={9}><Trans>Hebrew</Trans></option>
    <option value={10}><Trans>Kana</Trans></option>
    <option value={11}><Trans>Korean</Trans></option>
    <option value={12}><Trans>Thai</Trans></option>
    <option value={0}><Trans>Unspecified Standard</Trans></option>
    </>
  )
}

export default TransliterationMethodCode;