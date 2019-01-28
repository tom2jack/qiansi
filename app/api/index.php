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
        return 1;
    }

    public function miao()
    {
        return [1, ['data' => [111,2222]], '测试'];
    }
}