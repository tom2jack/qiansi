<?php
/**
 * 无关紧要莫名其妙的首页
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2018-11-27
 * Time: 09:54
 */

namespace app\api;


use zhimiao\Request;
use zhimiao\Utils;

class index
{
    public function index()
    {
        $cipher = 'AES-256-CBC';
        $ivlen = openssl_cipher_iv_length($cipher);
        $iv = bin2hex(openssl_random_pseudo_bytes($ivlen));
        $ivx = '51ae84ba12c3a6ab991a89070555bae8';
        $ivx = hex2bin($ivx);
        $s = '123456';
        $data = \app\Service\Utils::encrypt('hello wor中?', $s);
        $data = \app\Service\Utils::decrypt('tZbvIfAS4sKkmgVtldmrXL4/f+p2kq8dPXpQwFKypoA=', $s);
        var_dump($data);exit;
        $param = array_merge(Request::get(), Request::post());
        var_dump($param);exit;
        return [1, Utils::decrypt(Request::post('data'))];
    }

    public function miao()
    {
        $secret = '123456';
        $data = \app\Service\Utils::decrypt(Request::raw(), $secret);
        if ($data == 'zhimiao.api.check') {
            echo \app\Service\Utils::encrypt('succ', $secret);
            exit;
        }
    }
}