import React from 'react';
import { useLingui } from "@lingui/react";


const TransliterationMethodCode = () => {
  const { i18n } = useLingui();
  return (
    <>
    <option value={1}>{i18n._("Arabic (Arabic language")}</option>
    <option value={2}>{i18n._("Arabic (Persian language")}</option>
    <option value={3}>{i18n._("Armenian")}</option>
    <option value={4}>{i18n._("Cyrillic")}</option>
    <option value={5}>{i18n._("Devanagari & related Indic")}</option>
    <option value={6}>{i18n._("Georgian")}</option>
    <option value={7}>{i18n._("Greek")}</option>
    <option value={8}>{i18n._("Han (Hanzi, Kanji, Hanja")}</option>
    <option value={9}>{i18n._("Hebrew")}</option>
    <option value={10}>{i18n._("Kana")}</option>
    <option value={11}>{i18n._("Korean")}</option>
    <option value={12}>{i18n._("Thai")}</option>
    <option value={0}>{i18n._("Unspecified Standard")}</option>
    </>
  )
}

export default TransliterationMethodCode;