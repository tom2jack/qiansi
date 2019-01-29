<?php
/**
 *
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2018-11-27
 * Time: 09:54
 */

namespace app\api;


class index
{
    public function index()
    {
        $ivlen = openssl_cipher_iv_length('AES-256-CBC');
        $iv = openssl_random_pseudo_bytes($ivlen);
        echo $iv;
    }

    public function miao()
    {
        return [1, ['data' => [111,2222]], '测试'];
    }
}