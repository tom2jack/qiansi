<?php
/**
 * 工具类
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-01-29
 * Time: 12:02
 */
namespace app\Service;
use \zhimiao\Config;

class Utils {

    /**
     * 加密
     * @param {string} $encryptStr 要加密的字符串
     * @param {string} $localIV 加密的向量
     * @param {string} $encryptKey 加密的密钥
     * @return {string}
     */
    static public function encrypt($encryptStr)
    {
        return openssl_encrypt($encryptStr, Config::get('openssl_secret.method'), Config::get('openssl_secret.key'), 0, Config::get('openssl_secret.iv'));
    }

    /**
     * 解密
     * @param {string} $encryptStr 要解密的字符串
     * @param {string} $localIV 解密的向量
     * @param {string} $encryptKey 密钥
     * @return {string}
     */
    static public function decrypt($encryptStr)
    {
        return openssl_decrypt($encryptStr, Config::get('openssl_secret.method'), Config::get('openssl_secret.key'), 0, Config::get('openssl_secret.iv'));
    }
}