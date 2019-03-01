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
        $param = array_merge(Request::get(), Request::post());
        var_dump($param);exit;
        return [1, Utils::decrypt(Request::post('data'))];
    }

    public function miao()
    {
        return [1, ['data' => [111,2222]], '测试'];
    }
}